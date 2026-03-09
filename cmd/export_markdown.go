package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/riba2534/feishu-cli/internal/client"
	"github.com/riba2534/feishu-cli/internal/config"
	"github.com/riba2534/feishu-cli/internal/converter"
	"github.com/spf13/cobra"
)

var exportMarkdownCmd = &cobra.Command{
	Use:   "export <document_id|url>",
	Short: "导出文档为 Markdown",
	Long: `将飞书文档导出为 Markdown 格式。

支持通过文档 ID 或 URL 导出：
  feishu-cli doc export ABC123def456
  feishu-cli doc export https://xxx.feishu.cn/docx/ABC123def456
  feishu-cli doc export https://xxx.larkoffice.com/docx/ABC123def456

使用 --download-images 可同时下载文档中的图片和画板（画板自动导出为 PNG），
通过 --assets-dir 指定资源保存目录（默认 ./assets）。

示例:
  feishu-cli doc export ABC123def456
  feishu-cli doc export ABC123def456 --output doc.md
  feishu-cli doc export ABC123def456 --download-images
  feishu-cli doc export ABC123def456 --download-images --assets-dir ./images`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Validate(); err != nil {
			return err
		}

		documentID, err := extractDocToken(args[0])
		if err != nil {
			return err
		}
		output, _ := cmd.Flags().GetString("output")
		downloadImages, _ := cmd.Flags().GetBool("download-images")
		assetsDir, _ := cmd.Flags().GetString("assets-dir")

		// Get all blocks
		blocks, err := client.GetAllBlocks(documentID)
		if err != nil {
			return fmt.Errorf("获取块失败: %w", err)
		}

		frontMatter, _ := cmd.Flags().GetBool("front-matter")
		highlight, _ := cmd.Flags().GetBool("highlight")

		expandMentions, _ := cmd.Flags().GetBool("expand-mentions")

		// Convert to Markdown
		options := converter.ConvertOptions{
			DownloadImages: downloadImages,
			AssetsDir:      assetsDir,
			DocumentID:     documentID,
			FrontMatter:    frontMatter,
			Highlight:      highlight,
			ExpandMentions: expandMentions,
		}

		// 创建转换器（支持读取内嵌电子表格和展开 @用户）
		sheetReader := &FeishuSheetReader{}
		var userResolver converter.UserResolver
		if expandMentions {
			userResolver = &FeishuUserResolver{}
		}
		conv := converter.NewBlockToMarkdownFull(blocks, options, sheetReader, userResolver)
		markdown, err := conv.Convert()
		if err != nil {
			return fmt.Errorf("转换为 Markdown 失败: %w", err)
		}

		// 添加 Front Matter
		if frontMatter {
			docTitle := ""
			doc, docErr := client.GetDocument(documentID)
			if docErr == nil && doc != nil && doc.Title != nil {
				docTitle = *doc.Title
			}
			fm := fmt.Sprintf("---\ntitle: %q\ndocument_id: %s\n---\n\n", docTitle, documentID)
			markdown = fm + markdown
		}

		// Output
		if output != "" {
			if err := os.WriteFile(output, []byte(markdown), 0644); err != nil {
				return fmt.Errorf("写入输出文件失败: %w", err)
			}
			fmt.Printf("已导出到 %s\n", output)
		} else {
			fmt.Print(markdown)
		}

		return nil
	},
}

// extractDocToken 从 URL 或直接的 token 中提取 document_id
func extractDocToken(input string) (string, error) {
	// 尝试匹配 docx URL
	re := regexp.MustCompile(`/docx/([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(input)
	token := input
	if len(matches) > 1 {
		token = matches[1]
	}

	// 验证 token 格式
	if !isValidToken(token) {
		return "", fmt.Errorf("无效的文档 token: %s", token)
	}

	return token, nil
}

func init() {
	docCmd.AddCommand(exportMarkdownCmd)
	exportMarkdownCmd.Flags().StringP("output", "o", "", "输出文件路径")
	exportMarkdownCmd.Flags().Bool("download-images", false, "下载图片和画板到本地目录（画板自动导出为 PNG）")
	exportMarkdownCmd.Flags().String("assets-dir", "./assets", "图片和画板的保存目录")
	exportMarkdownCmd.Flags().Bool("front-matter", false, "添加 YAML front matter (标题和文档 ID)")
	exportMarkdownCmd.Flags().Bool("highlight", false, "保留文本颜色和背景色 (输出为 HTML span)")
	exportMarkdownCmd.Flags().Bool("expand-entions", true, "展开 @用户为友好格式 (需要 contact:user.base:readonly 权限)")
}

// FeishuSheetReader 实现 converter.SheetReader 接口
type FeishuSheetReader struct{}

// ReadSheet 读取内嵌电子表格内容
// sheetToken 格式为 "{spreadsheet_token}_{sheet_id}"
func (r *FeishuSheetReader) ReadSheet(sheetToken string) (*converter.SheetData, error) {
	// 解析 token：格式为 {spreadsheet_token}_{sheet_id}
	parts := strings.SplitN(sheetToken, "_", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid sheet token format: %s", sheetToken)
	}

	spreadsheetToken := parts[0]
	sheetID := parts[1]

	// 获取工作表信息以确定行列数
	sheets, err := client.QuerySheets(context.Background(), spreadsheetToken)
	if err != nil {
		return nil, fmt.Errorf("获取工作表列表失败: %w", err)
	}

	var targetSheet *client.SheetInfo
	for _, s := range sheets {
		if s.SheetID == sheetID {
			targetSheet = s
			break
		}
	}
	if targetSheet == nil {
		return nil, fmt.Errorf("未找到工作表: %s", sheetID)
	}

	// 读取表格内容
	rangeStr := fmt.Sprintf("%s!A1:%s%d", sheetID, columnLetter(targetSheet.ColCount), targetSheet.RowCount)
	cellRange, err := client.ReadCells(context.Background(), spreadsheetToken, rangeStr, "ToString", "FormattedString")
	if err != nil {
		return nil, fmt.Errorf("读取单元格失败: %w", err)
	}

	// 转换为 SheetData
	result := &converter.SheetData{
		Values: make([][]string, 0),
	}

	if cellRange != nil && cellRange.Values != nil {
		for _, row := range cellRange.Values {
			rowData := make([]string, 0)
			for _, cell := range row {
				cellStr := ""
				if cell != nil {
					switch v := cell.(type) {
					case string:
						cellStr = v
					case float64:
						cellStr = fmt.Sprintf("%.0f", v)
					case nil:
						cellStr = ""
					default:
						cellStr = fmt.Sprintf("%v", v)
					}
				}
				rowData = append(rowData, cellStr)
			}
			result.Values = append(result.Values, rowData)
		}
	}

	return result, nil
}

// columnLetter 将列号转换为 Excel 列字母（1=A, 2=B, ..., 27=AA）
func columnLetter(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+(col%26))) + result
		col /= 26
	}
	return result
}

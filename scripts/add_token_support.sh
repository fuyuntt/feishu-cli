#!/bin/bash
# 为 client API 调用添加 tokenOpt 支持的辅助脚本
# 用法: ./add_token_support.sh <file.go>

FILE=$1
if [ -z "$FILE" ]; then
    echo "Usage: $0 <file.go>"
    exit 1
fi

# 备份原文件
cp "$FILE" "${FILE}.bak"

# 使用 awk 来处理文件
awk '
BEGIN {
    in_func = 0
    need_token = 0
    brace_count = 0
}

# 检测函数开始
/func [A-Z]/ {
    in_func = 1
    brace_count = 0
    need_token = 0
}

# 检测 GetClient 调用
/GetClient\(\)/ {
    need_token = 1
}

# 在 GetClient 错误检查后插入 token 获取代码
need_token && /return.*err/ {
    print
    # 获取返回值类型
    getline
    print
    if ($0 ~ /^}$/ || $0 ~ /^$/) {
        # 插入 token 获取代码
        print ""
        print "\ttokenOpt, tokenErr := GetUserTokenOption()"
        print "\tif tokenErr != nil {"
        # 需要根据函数返回类型生成正确的 return 语句
        print "\t\treturn // TODO: 填充正确的返回值"
        print "\t}"
        print ""
    }
    need_token = 0
}

# 为 API 调用添加 tokenOpt 参数
/client\.[A-Za-z]+\.[A-Za-z]+\.[A-Za-z]+\(Context\(\),/ {
    # 在 ) 之前添加 tokenOpt
    gsub(/\)$/, ", tokenOpt)")
}

{ print }
' "${FILE}.bak" > "$FILE"

echo "Updated $FILE"
echo "Backup saved to ${FILE}.bak"

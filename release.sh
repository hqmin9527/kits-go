#!/bin/bash

# 获取最新 tag（按时间排序）
latest=$(git tag --sort=-v:refname | head -n 1)

if [ -z "$latest" ]; then
    echo "No tag found. Initializing to v0.1.0"
    new="v0.1.0"
else
    echo "Current latest tag: $latest"
    # 去掉 v
    ver=${latest#v}

    IFS='.' read -r major minor patch <<< "$ver"

    patch=$((patch + 1))

    new="v$major.$minor.$patch"
fi

echo "New version: $new"

# 创建标签
git tag "$new"
echo "Tag created: $new"

# 如果你想自动推送，取消注释：
git push origin "$new"

git commit -a

# 获取最后一个tag
last_tag=$(git tag | sort -V | tail -n 1)
#last_tag=$(git describe --abbrev=0 --tags)

# 提取版本号
version=$(echo $last_tag | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')

# 分割版本号
major=$(echo $version | cut -d. -f1)
minor=$(echo $version | cut -d. -f2)
patch=$(echo $version | cut -d. -f3)

# 递增版本号
patch=$((patch + 1))

# 生成新的tag
new_tag="v${major}.${minor}.${patch}"

echo "Last tag: $last_tag"
echo "New tag: $new_tag"

git tag $new_tag
git push origin $new_tag

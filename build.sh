echo "Welcome to use SMLKBOT!"
if [ "$1" == "TencentSCF" ]; then
    echo "Build SMLKBOT for Tencent SCF..."
    flags="-s -w -X 'SMLKBOT/utils/smlkshell.date=$(date)' -X 'SMLKBOT/utils/smlkshell.version="DevBuild-SCF"' -X 'SMLKBOT/utils/smlkshell.commit="$(git rev-parse --short HEAD)"' -X 'SMLKBOT/utils/smlkshell.IsSCF="${1}"'"
    go build -ldflags "$flags" -o SMLKBOTSCF
    echo "Build succeed!"
elif [ "$1" == "HTTP" ]; then
    echo "Build SMLKBOT for common"
    flags="-s -w -X 'SMLKBOT/utils/smlkshell.date=$(date)' -X 'SMLKBOT/utils/smlkshell.commit="$(git rev-parse --short HEAD)"'"
    go build -ldflags "$flags"
    echo "Build succeed!"
else
    echo "Oops, it seems you don't set build goal..."
    echo "Usage: bash build.sh [goal]"
    echo "Support goals:
    TencentSCF
    HTTP"
fi

echo "Welcome to use SMLKBOT!"
if [ "$1" == "TencentSCF" ]; then
    echo "Build SMLKBOT for Tencent SCF..."
    flags="-s -w -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.date=$(date)' -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.version="DevBuild-SCF"' -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.commit="$(git rev-parse --short HEAD)"' -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.IsSCF="${1}"'"
    go build -ldflags "$flags" -o ./target/SMLKBOTSCF
    echo "Build succeed!"
elif [ "$1" == "HTTP" ]; then
    echo "Build SMLKBOT for common"
    flags="-s -w -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.date=$(date)' -X 'github.com/Ink-33/SMLKBOT/utils/smlkshell.commit="$(git rev-parse --short HEAD)"'"
    if [ "$2" == "win" ]; then
        GOOS=windows go build -ldflags "$flags" -o ./target/SMLKBOT-win.exe
        echo "Build succeed!"
    elif [ "$2" == "arm" ]; then
        GOARCH=arm64 go build -ldflags "$flags" -o ./target/SMLKBOT-arm64
        echo "Build succeed!"
    else
        go build -ldflags "$flags" -o ./target/SMLKBOT
        echo "Build succeed!"
    fi
else
    echo "Oops, it seems you don't set build goal..."
    echo "Usage: bash build.sh [goal]"
    echo "Support goals:
    TencentSCF
    HTTP"
fi

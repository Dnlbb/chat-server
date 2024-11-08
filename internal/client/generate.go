package client

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/Dnlbb/platform_common/pkg/db.TxManager  -o ./mocks/ -s "_minimock.go"

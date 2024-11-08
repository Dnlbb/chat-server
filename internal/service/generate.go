package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ./servinterfaces.ChatService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i github.com/Dnlbb/auth/pkg/auth_v1.AuthClient  -o ./mocks/ -s "_minimock.go"

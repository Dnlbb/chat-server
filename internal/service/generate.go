package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ./servinterfaces.ChatService -o ./mocks/ -s "_minimock.go"

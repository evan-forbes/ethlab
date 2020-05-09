package server

// func echo(ctx context.Context, conn *websocket.Conn) {
// 	defer conn.Close()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			input := make(map[string]interface{})
// 			err := conn.ReadJSON(input)
// 			if err != nil {
// 				log.Printf("failure to read json during echo: %s", err)
// 			}
// 			fmt.Println(input)
// 		}
// 	}
// }

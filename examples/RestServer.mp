let login = fn(writer, request) {
	method = request.method
	if method = http.METHOD_POST {
		writer.writeHeader(200)
		writer.write(json.toJson({ "sessionId": "3d5bd2cA15ef0476" }))
	}
}

let logout = fn(writer, request) {
	method = request.method
	if method = http.METHOD_POST {
		writer.writeHeader(201)
	}
}

http.handleFunc("/authentication/login", login)
http.handleFunc("/authentication/logout", logout)
http.listenAndServe("127.0.0.1:8080")
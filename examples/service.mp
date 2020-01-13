service Hello on "0.0.0.0:8090" {
  @route(url="/authentication/login", methods=["POST"])
  fn login(writer, request) {
    writer.header.add("Content-Type", "application/json; charset=utf-8")
    writer.write(json.toJson({ "sessionId": "3d5bd2cA15ef047689" }))
  }

  @route(url="/authentication/logout", methods=["POST"])
  fn logout(writer, request) {
    writer.writeHeader(http.STATUS_CREATED) // return http status code 201
  }

  @route(url="/meters/setting-result/{acceptNo}", methods=["GET"])
  fn load_survey_result(writer, request) {
    writer.header.add("Content-Type", "application/json; charset=utf-8")
    //using 'vars' dictionary to access the url parameters
    writer.write(json.toJson({ "acceptNo": vars["acceptNo"], "resultCode": "1"}))
  }

  @route(url="/articles/{category}/{id:[0-9]+}", methods=["GET"])
  fn load_survey_result(writer, request) {
    writer.header.add("Content-Type", "application/json; charset=utf-8")
    //using 'vars' dictionary to access the url parameters
    writer.write(json.toJson({ "category": vars["category"], "id": vars["id"]}))
  }
}
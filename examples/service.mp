//service Hello on "0.0.0.0:8090" {
service Hello on "0.0.0.0:8090:debug" { //':debug': for debugging request
  @route(url="/authentication/login", methods=["POST"])
  fn login(writer, request) {
    writer.writeJson({ sessionId: "3d5bd2cA15ef047689" })
  }

  @route(url="/authentication/logout", methods=["POST"])
  fn logout(writer, request) {
    writer.writeHeader(http.STATUS_CREATED) // return http status code 201
  }

  @route(url="/meters/setting-result/{acceptNo}", methods=["GET"])
  fn load_survey_result(writer, request) {
    //using 'vars' dictionary to access the url parameters
    writer.writeJson({ acceptNo: vars["acceptNo"], resultCode: "1"})
  }

  @route(url="/articles/{category}/{id:[0-9]+}", methods=["GET"])
  fn load_survey_result(writer, request) {
    //using 'vars' dictionary to access the url parameters
    writer.writeJson({ category: vars["category"], id: vars["id"]})
  }
}
//service Hello on "0.0.0.0:8090" {
service Hello on "0.0.0.0:8090:debug" { //':debug': for debugging request
  @route(url="/authentication/login", methods=["POST"])
  fn login(writer, request) {
    //writer.writeJson({ sessionId: "3d5bd2cA15ef047689" })
    //writer.writeJson({ sessionId: "3d5bd2cA15ef047689" }), 200 # same as above
    //return { sessionId: "3d5bd2cA15ef047689" }, 200 # same as above
    return { sessionId: "3d5bd2cA15ef047689" } //same as above
  }

  @route(url="/authentication/logout", methods=["POST"])
  fn logout(writer, request) {
    // writer.writeHeader(http.STATUS_CREATED) # return http status code 201
    return http.STATUS_CREATED //same as above
  }

  @route(url="/meters/setting-result/{acceptNo}", methods=["GET"])
  fn load_survey_result(writer, request) {
    //using 'vars' dictionary to access the url parameters
    //writer.writeJson({ acceptNo: vars["acceptNo"], resultCode: "1"})
    return { acceptNo: vars["acceptNo"], resultCode: "1"} //same as above
  }

  @route(url="/articles/{category}/{id:[0-9]+}", methods=["GET"])
  fn getArticle(writer, request) {
    //using 'vars' dictionary to access the url parameters
    //writer.writeJson({ category: vars["category"], id: vars["id"]})
    return { category: vars["category"], id: vars["id"]} //same as above
  }
}

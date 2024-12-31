import "package:http/src/response.dart";
import 'package:http/http.dart' as http;
import 'dart:developer' as developer;
import 'dart:convert';

import '../globalData.dart';

class errorHandler {
  static Response apiRequestFailedMessage(Object error) {
    String result = '$error';
    if (result.split(' ').contains("SocketException:") == true ||
        (result.split(' ').contains("Invalid") == true &&
        result.split(' ').contains("argument(s):") == true)) {
      result = Globaldata.ErrorServerConnection;
    } else {
      result = Globaldata.ErrorBackendDown;
    }
    Response response = http.Response(result, 504);
    return response;
  }

  static Future<dynamic> basicResponseErrorHandler(Response response, String errorName) async {
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "$errorName FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    return responseData;
  }
}
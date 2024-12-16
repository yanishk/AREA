import "package:http/src/response.dart";
import 'package:http/http.dart' as http;

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
}
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import "package:http/src/response.dart";

import '../globalData.dart';

class GetApi {
  static Future<Response> get(String apiRoute) async {
    String apiUrl = dotenv.get('SERVER_LINK',
        fallback: 'Cannot find .env variable "SERVER_LINK"');
    String result = "";
    developer.log('GetApi route : $apiRoute');
    try {
      final response = await http.get(
        Uri.parse(apiUrl + apiRoute),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("GetApi responce :${jsonDecode(response.body)}\nstatusCode : ${response.statusCode}");
      return response;
    } catch (e) {
      developer.log("GetApi ERROR: $e", level: 1);
      result = '$e';
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
}

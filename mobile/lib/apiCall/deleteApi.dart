import 'package:http/http.dart' as http;
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import "package:http/src/response.dart";
import 'dart:convert';

import '../globalData.dart';

class DeleteApi {
  static Future<Response> delete(String apiRoute) async {
    String apiUrl = dotenv.get('SERVER_LINK',
        fallback: 'Cannot find .env variable "SERVER_LINK"');
    String result = "";
    developer.log('DeleteApi route : $apiRoute');
    try {
      final response = await http.delete(
        Uri.parse('${apiUrl + apiRoute}'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("DeleteApi responce :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("DeleteApi ERROR: $e");
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

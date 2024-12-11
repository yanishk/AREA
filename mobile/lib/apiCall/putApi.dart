import 'package:http/http.dart' as http;
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import "package:http/src/response.dart";
import 'dart:convert';

import '../globalData.dart';

class PutApi {
  static Future<Response> put(String apiRoute, Object? body) async {
    String apiUrl = dotenv.get('SERVER_LINK',
        fallback: 'Cannot find .env variable "SERVER_LINK"');
    String result = "";
    developer.log('PutApi route : $apiRoute');
    try {
      final response = await http.put(
        Uri.parse('${apiUrl + apiRoute}'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
        body: body
      );

      developer.log("PutApi responce :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("PutApi ERROR: $e");
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

import 'package:namer_app/apiCall/apiRequest.dart';
import 'dart:developer' as developer;
import 'dart:convert';

class ActionsOrReactionsFunctions {
  static Future<dynamic> getRouteResponse(String route) async {
    final response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getRouteResponse FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    developer.log("responseData : ${responseData}");
    return responseData;
  }

  static dynamic getMapKeyName(dynamic data) {
    String mapKeyName = data.toString().split(":").first.split("{").last;
    if (!(data[mapKeyName] is List<dynamic>)) {
      return getMapKeyName(data[mapKeyName]);
    }
    return data;
  }
}
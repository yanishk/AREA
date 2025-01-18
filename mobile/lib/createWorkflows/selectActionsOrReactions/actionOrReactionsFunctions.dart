import 'package:flutter/material.dart';
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
    return responseData;
  }

  static dynamic getMapKeyName(dynamic data) {
    String mapKeyName = data.toString().split(":").first.split("{").last;
    if (data == null) {
      return data;
    }
    if (!(data[mapKeyName] is List<dynamic>)) {
      return getMapKeyName(data[mapKeyName]);
    }
    return data;
  }

  static Color fadeColor(Color color, double ratio) {
    int r = color.red ~/ 2;
    int g = color.green ~/ 2;
    int b = color.blue ~/ 2;
    return Color.fromARGB(255, r, g, b);
  }
}
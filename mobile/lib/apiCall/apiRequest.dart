import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import "package:http/src/response.dart";

import "apiErrorHandling.dart";
import "apiData.dart";
import '../globalData.dart';

class ApiRequest {
  static Future<Response> get(String apiRoute) async {
    developer.log('GetApi route : $apiRoute');
    try {
      final response = await http.get(
        Uri.parse(ApiData.apiUrl + apiRoute),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("GetApi responce :${jsonDecode(response.body)}\nstatusCode : ${response.statusCode}");
      return response;
    } catch (e) {
      developer.log("GetApi ERROR: $e", level: 1);
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<Response> put(String apiRoute, Object? body) async {
    developer.log('PutApi route : $apiRoute');
    try {
      final response = await http.put(
        Uri.parse('${ApiData.apiUrl + apiRoute}'),
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
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<Response> delete(String apiRoute) async {
    developer.log('DeleteApi route : $apiRoute');
    try {
      final response = await http.delete(
        Uri.parse('${ApiData.apiUrl + apiRoute}'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("DeleteApi responce :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("DeleteApi ERROR: $e");
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<Response> post(String apiRoute, Object? body) async {
    developer.log('PostApi route : $apiRoute');
    try {
      final response = await http.post(
        Uri.parse('${ApiData.apiUrl + apiRoute}'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
        body: body
      );

      developer.log("PostApi responce :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("PostApi ERROR: $e");
      return errorHandler.apiRequestFailedMessage(e);
    }
  }
}

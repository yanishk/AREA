import 'package:http/http.dart' as http;
import 'package:http/io_client.dart';
import 'dart:convert';
import 'dart:developer' as developer;
import "package:http/src/response.dart";
import 'dart:io';

import "apiErrorHandling.dart";
import "apiData.dart";
import '../globalData.dart';

HttpClient createHttpClient() {
  var httpClient = HttpClient()
    ..badCertificateCallback = ((X509Certificate cert, String host, int port) => true);
  return httpClient;
}

class ApiRequest {
  static final ioClient = IOClient(createHttpClient());

  static String removeDoubleSlash(String url) {
    return "https://" + url.split("https://").last.replaceAll("//", "/");
  }

  static Future<Response> get(String apiRoute) async {
    developer.log('GetApi route : $apiRoute');
    try {
      final response = await ioClient.get(
        Uri.parse(removeDoubleSlash(ApiData.apiUrl + apiRoute)),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("GetApi response :${jsonDecode(response.body)}\nstatusCode : ${response.statusCode}");
      return response;
    } catch (e) {
      developer.log("GetApi ERROR: $e", level: 1);
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<Response> put(String apiRoute, Object? body) async {
    developer.log('PutApi route : $apiRoute');
    try {
      final response = await ioClient.put(
        Uri.parse(removeDoubleSlash(ApiData.apiUrl + apiRoute)),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
        body: body
      );

      developer.log("PutApi response :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("PutApi ERROR: $e");
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<Response> delete(String apiRoute) async {
    developer.log('DeleteApi route : $apiRoute');
    try {
      final response = await ioClient.delete(
        Uri.parse(removeDoubleSlash(ApiData.apiUrl + apiRoute)),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      developer.log("DeleteApi response :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("DeleteApi ERROR: $e");
      return errorHandler.apiRequestFailedMessage(e);
    }
  }

  static Future<http.Response> post(String apiRoute, Object? body) async {
    developer.log('PostApi route : ${ApiData.apiUrl}$apiRoute');
    developer.log("api domaine : ${ApiData.apiUrl}");
    try {
      final response = await ioClient.post(
        Uri.parse(removeDoubleSlash(ApiData.apiUrl + apiRoute)),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
        body: body
      );

      developer.log("PostApi response :${jsonDecode(response.body)}");
      return response;
    } catch (e) {
      developer.log("PostApi ERROR: $e");
      return errorHandler.apiRequestFailedMessage(e);
    }
  }
}

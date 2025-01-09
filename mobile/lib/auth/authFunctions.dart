import "package:http/src/response.dart";
import 'dart:developer' as developer;
import 'dart:convert';

import "../../globalData.dart";
import "../../apiCall/apiRequest.dart";
import 'package:flutter_inappwebview/flutter_inappwebview.dart';

class AuthFunctions {
  static Future<String> getAuthLink(String callBackType, String service) async {
    String authLink = "";
    Response response = await ApiRequest.get("authentication?service=$service&callbacktype=$callBackType&apptype=${Globaldata.appType}");
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return authLink;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getAuthLink, service : $service, FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return authLink;
    }
    authLink = responseData["auth-url"];
    developer.log(
      "getAuthLink, service : $service, response : ${authLink}");
    return authLink;
  }

  static Future<void> setCookiesAndNavigate(String apiUrl) async {
    if (Globaldata.JWToken.isEmpty) {
      return;
    }
    final cookieManager = CookieManager();
    await cookieManager.setCookie(
      url: WebUri(apiUrl),
      name: "JWToken",
      value: Globaldata.JWToken.split('=').last,
      isSecure: true,
      isHttpOnly: true,
    );
  }

  static String removeServiceFromUrl(String url, String serviceName) {
    final splitUrl = url.split("/$serviceName");
    String finalUrl = "${splitUrl.first}${splitUrl.last}";
    return finalUrl;
  }
}
import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:namer_app/apiCall/apiRequest.dart';
import 'dart:developer' as developer;
import 'dart:convert';
import "package:http/src/response.dart";
import "package:namer_app/auth/authFunctions.dart";

import "../../apiCall/apiErrorHandling.dart";
import "../../globalData.dart";
import "../../login/loginNavigator.dart";
import "../../login/loadServicesList.dart";
import "../../apiCall/apiData.dart";
import "../authData.dart";

String authConnectionType = "";
String authServiceName = "";

class WebviewAuthPage extends StatelessWidget {
  WebviewAuthPage({
    super.key,
    required this.callBackType,
    required this.authService,
  });
  final String callBackType;
  final String authService;

  @override
  Widget build(BuildContext context) {
    authConnectionType = callBackType;
    authServiceName = authService;
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: WebviewAuthAPICall(),
    );
  }
}

class WebviewAuthAPICall extends StatefulWidget {
  @override
  WebviewAuth createState() => WebviewAuth();
}

class WebviewAuth extends State<WebviewAuthAPICall> {
  late InAppWebViewController webViewController;
  String newUrl = "";
  String apiRoute = "service-callback/$authServiceName";

  @override
  void initState() {
    super.initState();
    if (authConnectionType == "login") {
      developer.log("authConnectionType : $authConnectionType");
      setState(() {
        apiRoute = "login-callback/$authServiceName";
      });
    }
  }

  Future<void> saveJWToken(Response response) async {
    String? cookiesBrute = response.headers["set-cookie"];
    String JWToken = "";
    if (cookiesBrute != null) {
      JWToken = cookiesBrute.split(';').first;
      setState(() {
        developer.log("JWToken : {${JWToken}}");
        Globaldata.JWToken = JWToken;
        LoginNavigator.redirectAndHandleGoBackRedirection(LoadServicePage(), LoadServicePage());
      });
      return;
    }
  }

  String removeUrlState(String url) {
    final splitState = url.split("?state=");
    if (splitState.length == 1) {
      return url;
    }
    developer.log("removeUrlState splitState : $splitState");
    final splitCode = splitState.last.split("&code=");
    developer.log("removeUrlState splitCode : $splitCode");
    return "${splitState.first}?code=${splitCode.last}";
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: InAppWebView(
        initialUrlRequest: URLRequest(
          url: WebUri(AuthData.serviceConnectionLink),
        ),
        initialSettings: InAppWebViewSettings(
          javaScriptEnabled: true,
          domStorageEnabled: true,
        ),
        onLoadStart: (controller, url) {
          AuthFunctions.setCookiesAndNavigate(ApiData.apiUrl);
        },
        onLoadStop: (controller, url) async {
          developer.log("onLoadStop url : ${url.toString()}");
          if (url != null &&
              url.toString().contains(ApiData.apiUrl + apiRoute)) {
            final urlSplit = url.toString().split(ApiData.apiUrl);
            String callBackUrl = removeUrlState(urlSplit.last);
            final Response response = await ApiRequest.post(AuthFunctions.insertAppTypeInUrl(callBackUrl), jsonEncode(<String, dynamic>{}),);

            if (response.statusCode == 200) {
              developer.log("User connected");
              if (Globaldata.JWToken.isEmpty) {
                saveJWToken(response);
              }
              Navigator.pop(context);
              Navigator.pop(Globaldata.myContext);
            } else {
              final responseData = await errorHandler.basicResponseErrorHandler(response, "Error auth failed : ");
              if (responseData == null) {
                return;
              }
            }
          }
        },
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'dart:developer' as developer;
import "package:namer_app/auth/authFunctions.dart";

import "../../globalData.dart";
import "../../apiCall/apiData.dart";
import "../authData.dart";
import 'AuthLoadingScreen.dart';

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
  String apiRoute = "service-callback";
  String callBackRedirect = "create";

  @override
  void initState() {
    super.initState();
    if (authConnectionType == "login") {
      developer.log("authConnectionType : $authConnectionType");
      setState(() {
        apiRoute = "login-callback";
        callBackRedirect = "explore";
      });
    }
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
        onLoadStart: (controller, url) async {
          AuthFunctions.setCookiesAndNavigate(ApiData.apiUrl);
          developer.log("onLoadStart url : $url");
          if (url != null && url.toString().contains("/$callBackRedirect?") ||
              url != null && url.toString().contains("/$authServiceName?")) {
            final authSucced = await Navigator.push(
              Globaldata.myContext,
              MaterialPageRoute(builder: (context) => AuthLoadingScreenPage(
                callBack: callBackRedirect,
                authService: authServiceName,
                route: apiRoute,
                uri: url,
              )),
            );
            Navigator.pop(context);
            Navigator.pop(Globaldata.myContext, authSucced);
            return;
          }
        },
      ),
    );
  }
}

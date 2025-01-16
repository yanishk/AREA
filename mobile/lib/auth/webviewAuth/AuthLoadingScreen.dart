import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:loading_animation_widget/loading_animation_widget.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:namer_app/apiCall/apiRequest.dart';

import "../../apiCall/apiErrorHandling.dart";
import "../../globalData.dart";
import "../../login/loginNavigator.dart";
import "../../login/loadData/loadServicesData.dart";

String callBackRedirect = "";
String authServiceName = "";
String apiRoute = "";
WebUri webUri = WebUri("");

class AuthLoadingScreenPage extends StatelessWidget {
  AuthLoadingScreenPage({
    super.key,
    required this.callBack,
    required this.authService,
    required this.route,
    required this.uri,
  });
  final String callBack;
  final String authService;
  final String route;
  final WebUri uri;
  @override
  Widget build(BuildContext context) {
    callBackRedirect = callBack;
    authServiceName = authService;
    apiRoute = route;
    webUri = uri;
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  AuthLoadingScreen createState() => AuthLoadingScreen();
}

class AuthLoadingScreen extends State<APICall> {
  @override
  void initState() {
    super.initState();
    waitAuthEnd();
  }

  Future<void> waitAuthEnd() async {
    await handleCallbackUrl();
  }

  Future<void> saveJWToken(Response response) async {
    String? cookie = response.headers["set-cookie"];
    String JWToken = "";
    if (cookie != null) {
      JWToken = cookie.split(';').first;
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
      return url.split("&state=").first;
    }
    final splitCode = splitState.last.split("&code=");
    return "${splitState.first}?code=${splitCode.last}";
  }

  Future<void> handleCallbackUrl() async {
    List<String> urlSplit = webUri.toString().split(callBackRedirect);
    if (webUri.toString().contains("/$authServiceName?")) {
      urlSplit = webUri.toString().split(authServiceName);
    }
    String callBackUrl = removeUrlState(urlSplit.last);
    String postUrl = "${apiRoute}$callBackUrl";
    final Response response = await ApiRequest.post(postUrl,
      jsonEncode(<String, dynamic>{
        "apptype" : Globaldata.appType,
        "service" : authServiceName,
      }),
    );

    if (response.statusCode == 200) {
      developer.log("User connected");
      if (Globaldata.JWToken.isEmpty) {
        saveJWToken(response);
      }
      Navigator.pop(Globaldata.myContext, true);
    } else {
      errorHandler.basicResponseErrorHandler(response, "Error auth failed : ");
      Navigator.pop(Globaldata.myContext);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: LoadingAnimationWidget.dotsTriangle(
          color: Color.fromARGB(255, 34, 34, 34),
          size: 200,
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
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

String redditAuthType = "";

class RedditAuthPage extends StatelessWidget {
  RedditAuthPage({
    super.key,
    required this.callBackType,
  });
  final String callBackType;
  @override
  Widget build(BuildContext context) {
    redditAuthType = callBackType;
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: RedditAuthAPICall(),
    );
  }
}

class RedditAuthAPICall extends StatefulWidget {
  @override
  RedditAuth createState() => RedditAuth();
}

class RedditAuth extends State<RedditAuthAPICall> {
  late InAppWebViewController webViewController;
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  String newUrl = "";
  String apiRoute = "service-callback/Reddit";

  @override
  void initState() {
    super.initState();
    if (redditAuthType == "login") {
      developer.log("redditAuthType : $redditAuthType");
      setState(() {
        apiRoute = "login-callback/Reddit";
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
          url: WebUri(AuthData.redditConnectionLink),
        ),
        initialSettings: InAppWebViewSettings(
          javaScriptEnabled: true,
          domStorageEnabled: true,
        ),
        onLoadStart: (controller, url) {
          AuthFunctions.setCookiesAndNavigate(apiUrl);
        },
        onLoadStop: (controller, url) async {
          developer.log("onLoadStop url : ${url.toString()}");
          if (url != null &&
              url.toString().contains(apiUrl + apiRoute)) {
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

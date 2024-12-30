import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:namer_app/apiCall/apiRequest.dart';
import 'dart:developer' as developer;
import 'dart:convert';
import "package:http/src/response.dart";

import "../apiCall/apiErrorHandling.dart";
import "../globalData.dart";
import "../login/loginNavigator.dart";
import "../login/loadServicesList.dart";

String spotifyAuthType = "";

class SpotifyAuthPage extends StatelessWidget {
  SpotifyAuthPage({
    super.key,
    required this.callBackType,
  });
  final String callBackType;
  @override
  Widget build(BuildContext context) {
    spotifyAuthType = callBackType;
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: SpotifyAuthAPICall(),
    );
  }
}

class SpotifyAuthAPICall extends StatefulWidget {
  @override
  SpotifyAuth createState() => SpotifyAuth();
}

class SpotifyAuth extends State<SpotifyAuthAPICall> {
  late InAppWebViewController webViewController;
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  String newUrl = "";
  String apiRoute = "service-callback/Spotify";

  @override
  void initState() {
    super.initState();
    if (spotifyAuthType == "login") {
      developer.log("spotifyAuthType : $spotifyAuthType");
      setState(() {
        apiRoute = "login-callback/Spotify";
      });
    }
  }

  Future<void> _setCookiesAndNavigate() async {
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

  String insertAppTypeInUrl(String url) {
    final splitUrl = url.split("?code=");
    String finalUrl = "${splitUrl.first}/${Globaldata.appType}?code=${splitUrl.last}";
    return finalUrl;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: InAppWebView(
        initialUrlRequest: URLRequest(
          url: WebUri(Globaldata.spotifyConnectionLink),
        ),
        initialSettings: InAppWebViewSettings(
          javaScriptEnabled: true,
          domStorageEnabled: true,
        ),
        onLoadStart: (controller, url) {
          _setCookiesAndNavigate();
        },
        onLoadStop: (controller, url) async {
          developer.log(url.toString());
          if (url != null &&
              url.toString().contains(apiUrl + apiRoute)) {
            final urlSplit = url.toString().split(':8080/');
            final Response response = await ApiRequest.post(insertAppTypeInUrl(urlSplit.last), jsonEncode(<String, dynamic>{}),);

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

import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import "globalData.dart";

class SpotifyAuthPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
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
  @override
  void initState() {
    super.initState();
    _setCookiesAndNavigate();
  }
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');

  Future<void> _setCookiesAndNavigate() async {
    final cookieManager = CookieManager();
    await cookieManager.setCookie(
      url: WebUri(apiUrl),
      name: "JWToken",
      value: Globaldata.JWToken.split('=').last,
      isSecure: true,
      isHttpOnly: true,
    );
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
        onWebViewCreated: (controller) {
          webViewController = controller;
        },
      ),
    );
  }
}

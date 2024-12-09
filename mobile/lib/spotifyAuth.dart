import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'dart:developer' as developer;

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
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  String newUrl = "";

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
        onLoadStart: (controller, url) {
          _setCookiesAndNavigate();
        },
        onLoadStop: (controller, url) async {
          if (url != null && url.toString().contains("${apiUrl}callback/Spotify")) {
            final content = await controller.evaluateJavascript(source: "document.body.innerText");
            if (content != null && content.contains('"success": "Token generated"')) {
              developer.log("User connected :\n$content");
              Navigator.pop(context);
              Navigator.pop(Globaldata.myContext);
            } else {
              developer.log("Error auth failed :\n$content");
            }
          }
        },
      ),
    );
  }
}

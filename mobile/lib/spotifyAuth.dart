import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_web_auth/flutter_web_auth.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import 'package:url_launcher/url_launcher.dart';

import "custom_widget.dart";
import "globalData.dart";
import 'getApi.dart';

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
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  List<String> serviceNameList = Globaldata.actionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons;
  Color colorTheme = Colors.black;
  String title = "Select action";
  String serviceName = "";
  String result = '';
  String apiRoute = "callback/";
  String actionDescription = "";
  bool isInit = false;

  Future<void> setReactionService() async {
    // final result = await FlutterWebAuth2.authenticate(url: Globaldata.spotifyConnectionLink, callbackUrlScheme: "${apiUrl + apiRoute}Spotify");
    await launchUrl(Uri.parse(Globaldata.spotifyConnectionLink));

    final token = Uri.parse(result).queryParameters['token'];
    developer.log("token : $token");
    Navigator.pop(Globaldata.myContext);
  }

  @override
  Widget build(BuildContext context) {
    if (!isInit) {
      isInit = true;
      setReactionService();
    }
    return Column(
      children: const [
        MyDivider(),
      ],
    );
  }
}

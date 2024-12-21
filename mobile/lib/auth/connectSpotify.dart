import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'dart:developer' as developer;
import 'dart:convert';

import "../globalData.dart";
import "../apiCall/apiRequest.dart";
import "spotifyLoginWebview.dart";

class MySpotifyAuth {
  static Future<String> getSpotifyAuthLink(String callBackType) async {
    String spotifyLink = "";
    Response response = await ApiRequest.get("authentication/Spotify/$callBackType");
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return spotifyLink;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getSpotifyAuthLink FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return spotifyLink;
    }
    spotifyLink = responseData["auth-url"];
    developer.log(
      "getSpotifyAuthLink response : ${spotifyLink}");
    return spotifyLink;
  }

  static Future<void> connectSpotify(String callBackType) async {
    String spotifyLink = await getSpotifyAuthLink(callBackType);
    if (spotifyLink == "") {
      return;
    }
    Globaldata.spotifyConnectionLink = spotifyLink;
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => SpotifyAuthPage(callBackType: callBackType)),
    );
  }
}

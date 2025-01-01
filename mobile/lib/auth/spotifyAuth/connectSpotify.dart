import 'package:flutter/material.dart';

import "../../globalData.dart";
import "spotifyLoginWebview.dart";
import "../authFunctions.dart";
import "../authData.dart";

class MySpotifyAuth {
  static Future<void> connectSpotify(String callBackType) async {
    String spotifyLink = await AuthFunctions.getAuthLink(callBackType, "Spotify");
    if (spotifyLink == "") {
      return;
    }
    AuthData.spotifyConnectionLink = spotifyLink;
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => SpotifyAuthPage(callBackType: callBackType)),
    );
  }
}

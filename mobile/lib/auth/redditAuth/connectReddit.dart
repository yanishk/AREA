import 'package:flutter/material.dart';

import "../../globalData.dart";
import "redditLoginWebview.dart";
import "../authData.dart";
import "../authFunctions.dart";

class MyRedditAuth {
  static Future<void> connectReddit(String callBackType) async {
    String redditLink = await AuthFunctions.getAuthLink(callBackType, "Reddit");
    if (redditLink == "") {
      return;
    }
    AuthData.redditConnectionLink = redditLink;
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => RedditAuthPage(callBackType: callBackType)),
    );
  }
}

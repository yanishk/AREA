import 'package:flutter/material.dart';
import 'dart:developer' as developer;

import "../globalData.dart";

class LoginNavigator {
  static Future<void> redirectAndHandleGoBackRedirection(Widget page, Widget goBackPage) async {
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => page),
    );
    developer.log("Go back to unauthorized page");
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => goBackPage),
    );
  }
}
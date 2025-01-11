import 'package:flutter/material.dart';

import "../../../globalData.dart";
import "../authData.dart";
import "../authFunctions.dart";
import "AuthWebview.dart";

class ConnectServiceAuthWebview {
  static Future<bool> connectService(String callBackType, String serviceName) async {
    String authLink = await AuthFunctions.getAuthLink(callBackType, serviceName);
    if (authLink == "") {
      return false;
    }
    AuthData.serviceConnectionLink = authLink;
    final result = await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => WebviewAuthPage(callBackType: callBackType, authService: serviceName,)),
    );
    if (result != true) {
      return false;
    }
    return true;
  }
}

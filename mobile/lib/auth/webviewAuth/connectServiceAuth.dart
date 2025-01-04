import 'package:flutter/material.dart';

import "../../../globalData.dart";
import "../authData.dart";
import "../authFunctions.dart";
import "AuthWebview.dart";

class ConnectServiceAuthWebview {
  static Future<void> connectService(String callBackType, String serviceName) async {
    String authLink = await AuthFunctions.getAuthLink(callBackType, serviceName);
    if (authLink == "") {
      return;
    }
    AuthData.serviceConnectionLink = authLink;
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => WebviewAuthPage(callBackType: callBackType, authService: serviceName,)),
    );
  }
}

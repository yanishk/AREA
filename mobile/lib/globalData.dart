import 'package:flutter/material.dart';
import 'login_page.dart';

class Globaldata {
  static String JWToken = "";
  static String ErrorServerConnection =
      "Could not find the server\nPlease try again later";
  static BuildContext myContext = myContext;
  static Widget lastPage = LoginPage();
  static int firstService = -1;
  static int secondService = -1;

  static List<String> actionServicesName = [];
  static List<String> actionServicesColor = [];
  static List<String> actionServicesIcons = [];

  static List<String> reactionServicesName = [];
  static List<String> reactionServicesColor = [];
  static List<String> reactionServicesIcons = [];

  static String domainName = "https://mathisbrehm.fr/";
  static String spotifyConnectionLink = "";

  static ActionOrReactionData serviceAction = ActionOrReactionData.optional();
  static ActionOrReactionData serviceReaction = ActionOrReactionData.optional();

  static int getActionServicesSize() {
    return actionServicesName.length;
  }
}

class ActionOrReactionData {
  String id = "";
  String name = "";
  String description = "";
  String serviceId = "";
  int nbParam = 0;
  String parameter = "";
  Map<String, dynamic> parameterList = {};

  ActionOrReactionData.optional([
    this.id = "",
    this.name = "",
    this.description = "",
    this.serviceId = "",
    this.nbParam = 0,
    this.parameter = "",
    Map<String, dynamic>? parameterList,
  ]) {
    this.parameterList = parameterList ?? {};
  }
}

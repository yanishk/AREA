import 'package:flutter/material.dart';

class Globaldata {
  static String JWToken = "";
  static String ErrorServerConnection =
      "Could not find the server\nPlease try again later";
  static String ErrorBackendDown = "Gateway Time-out\nPlease try again later";
  static BuildContext myContext = myContext;
  static int firstService = -1;
  static int secondService = -1;

  static List<String> actionServicesName = [];
  static List<String> actionServicesColor = [];
  static List<String> actionServicesIcons = [];
  static List<String> actionServicesID = [];
  static List<ActionOrReactionList> actionServicesActions = [];

  static List<String> reactionServicesName = [];
  static List<String> reactionServicesColor = [];
  static List<String> reactionServicesIcons = [];
  static List<String> reactionServicesID = [];
  static List<ActionOrReactionList> reactionServicesReaction = [];

  static String domainName = "https://mathisbrehm.fr/";
  static String spotifyConnectionLink = "";

  static const String appType = "mobile";

  static void resetActionReactionServicesList() {
    Globaldata.actionServicesName = [];
    Globaldata.actionServicesColor = [];
    Globaldata.actionServicesIcons = [];
    Globaldata.actionServicesID = [];
    Globaldata.actionServicesActions = [];

    Globaldata.reactionServicesName = [];
    Globaldata.reactionServicesColor = [];
    Globaldata.reactionServicesIcons = [];
    Globaldata.reactionServicesID = [];
    Globaldata.reactionServicesReaction = [];
  }

  static int getActionServicesSize() {
    return actionServicesName.length;
  }
}

class ActionOrReactionList {
  String description = "";
  String id = "";
  String name = "";
  int nbparam = 0;
  String serviceid = "";
  String servicename = "";

  ActionOrReactionList.optional([
    this.description = "",
    this.id = "",
    this.name = "",
    this.nbparam = 0,
    this.serviceid = "",
    this.servicename = "",
  ]);
}
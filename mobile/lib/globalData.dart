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

  static List<String> reactionServicesName = [];
  static List<String> reactionServicesColor = [];
  static List<String> reactionServicesIcons = [];

  static String domainName = "https://mathisbrehm.fr/";
  static String spotifyConnectionLink = "";

  static void resetActionReactionServicesList() {
    Globaldata.actionServicesName = [];
    Globaldata.actionServicesColor = [];
    Globaldata.actionServicesIcons = [];
    Globaldata.reactionServicesName = [];
    Globaldata.reactionServicesColor = [];
    Globaldata.reactionServicesIcons = [];
  }

  static int getActionServicesSize() {
    return actionServicesName.length;
  }
}

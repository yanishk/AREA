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
  static List<bool> actionIsAuthNeeded = [];
  static List<ActionOrReactionList> actionServicesActions = [];

  static List<String> reactionServicesName = [];
  static List<String> reactionServicesColor = [];
  static List<String> reactionServicesIcons = [];
  static List<String> reactionServicesID = [];
  static List<bool> reactionIsAuthNeeded = [];
  static List<ActionOrReactionList> reactionServicesReaction = [];

  static List<ServiceData> serviceList = [];

  static String domainName = "https://mathisbrehm.fr/";

  static UserInfos userInfos = UserInfos.optional();

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
    Globaldata.userInfos = UserInfos.optional();

    Globaldata.serviceList = [];
  }
}

class ServiceData {
  String name;
  String color;
  String icon;
  String id;
  bool isAuthNeeded;
  bool hasActions;
  bool hasReactions;

  ServiceData.optional([
    this.name = "",
    this.color = "",
    this.icon = "",
    this.id = "",
    this.isAuthNeeded = false,
    this.hasActions = false,
    this.hasReactions = false,
  ]);
}

class ActionOrReactionList {
  String description = "";
  String id = "";
  String name = "";
  String serviceid = "";
  String servicename = "";
  List<ActionOrReactionsParameter> parameter = [];

  ActionOrReactionList.optional([
    this.description = "",
    this.id = "",
    this.name = "",
    this.serviceid = "",
    this.servicename = "",
    List<ActionOrReactionsParameter> ?parameter,
  ]) {
    this.parameter = parameter ?? [];
  }
}

class ActionOrReactionsParameter {
  bool isexhaustive;
  String name = "";
  String route = "";
  String type = "";
  List<String> values = [];

  ActionOrReactionsParameter.optional([
    this.isexhaustive = false,
    this.name = "",
    this.route = "",
    this.type = "",
    this.values = const [],
  ]);
}

class UserInfos {
  String email = "";
  String createdat = "";
  String timezone = "";
  String connectiontype = "";

  UserInfos.optional([
    this.email = "",
    this.createdat = "",
    this.timezone = "",
    this.connectiontype = "",
  ]);
}

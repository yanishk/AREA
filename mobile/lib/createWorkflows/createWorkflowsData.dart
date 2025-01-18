import "../globalData.dart";

bool isAnAction = true;

class CreateData {
  static ActionOrReactionData serviceAction = ActionOrReactionData.optional();
  static ActionOrReactionData serviceReaction = ActionOrReactionData.optional();
  static List<ServiceData> actualServiceList = [];
  static List<ServiceData> actionServiceList = [];
  static List<ServiceData> reactionServiceList = [];

  static ActionOrReactionData selectedActionOrReactionData = ActionOrReactionData.optional();
  static int selectedServiceId = 0;
  static bool selectOrUpdateWorkflow = false;

  static bool hasToLogin = false;

  static List<String> actionParamListName = [];
  static List<String> reactionParamListName = [];
}

class ActionOrReactionData {
  String description;
  String id;
  String name;
  String serviceId;
  String icon;
  List<ActionOrReactionsParameter> parameter = [];
  Map<String, dynamic> parameterList = {};

  ActionOrReactionData.optional([
    this.id = "",
    this.name = "",
    this.description = "",
    this.serviceId = "",
    this.icon = "",
    List<ActionOrReactionsParameter> ?parameter,
    Map<String, dynamic>? parameterList,
  ]) {
    this.parameterList = parameterList ?? {};
    this.parameter = parameter ?? [];
  }
}

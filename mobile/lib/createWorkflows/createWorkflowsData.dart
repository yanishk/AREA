import "../globalData.dart";

bool isAnAction = true;

class CreateData {
  static List<String> serviceNameList = Globaldata.actionServicesName;
  static List<String> serviceColorList = Globaldata.actionServicesColor;
  static List<String> serviceIconsList = Globaldata.actionServicesIcons;

  static ActionOrReactionData serviceAction = ActionOrReactionData.optional();
  static ActionOrReactionData serviceReaction = ActionOrReactionData.optional();
}

class ActionOrReactionData {
  String id;
  String name;
  String description;
  String serviceId;
  int nbParam;
  String parameter;
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

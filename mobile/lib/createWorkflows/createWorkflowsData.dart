import "../globalData.dart";

bool isAnAction = true;

class CreateData {
  static List<String> serviceNameList = Globaldata.actionServicesName;
  static List<String> serviceColorList = Globaldata.actionServicesColor;
  static List<String> serviceIconsList = Globaldata.actionServicesIcons;

  static ActionOrReactionData serviceAction = ActionOrReactionData.optional();
  static ActionOrReactionData serviceReaction = ActionOrReactionData.optional();

  static ActionOrReactionData selectedActionOrReactionData = ActionOrReactionData.optional();
  static int selectedServiceId = 0;
}

class ActionOrReactionData {
  String description;
  String id;
  String name;
  String serviceId;
  List<ActionOrReactionsParameter> parameter = [];
  Map<String, dynamic> parameterList = {};

  ActionOrReactionData.optional([
    this.id = "",
    this.name = "",
    this.description = "",
    this.serviceId = "",
    List<ActionOrReactionsParameter> ?parameter,
    Map<String, dynamic>? parameterList,
  ]) {
    this.parameterList = parameterList ?? {};
    this.parameter = parameter ?? [];
  }
}

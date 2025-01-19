import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;

import "../../globalData.dart";
import "../../explore/explorePage.dart";
import "../../apiCall/apiRequest.dart";
import "../../apiCall/apiErrorHandling.dart";
import "../loginNavigator.dart";
import 'loadServicesData.dart';

class LoadDataFunctions {
  static String getInfosFromUserDataString(String userData, String infos) {
    return (userData.split(infos).last).split(",").first;
  }

  static Future<void> storeUserInfos() async {
    Response response = await ApiRequest.get("user");
    final responseData = await errorHandler.basicResponseErrorHandler(response, "getEmail");
    if (responseData == null) {
      return;
    }
    final userDataString = (responseData.toString().split("{").last).split("}").first;
    Globaldata.userInfos.email = getInfosFromUserDataString(userDataString, "email: ");
    Globaldata.userInfos.createdat = getInfosFromUserDataString(userDataString, "createdat: ");
    Globaldata.userInfos.timezone = getInfosFromUserDataString(userDataString, "timezone: ");
    Globaldata.userInfos.connectiontype = getInfosFromUserDataString(userDataString, "connectiontype: ");
    Globaldata.userInfos.username = Globaldata.userInfos.email.split("@").first.replaceAll(".", "");
  }

  static Future<dynamic> getServiceActionsOrReactions(String route) async {
    Response response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getServiceActionsOrReactions FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    return responseData;
  }

  static List<String> getValueListe(List<dynamic> values) {
    List<String> valuesList = [];
    for (int i = 0; i < values.length; i++) {
      if (values[i] != null) {
        valuesList.add(values[i]);
      }
    }
    return valuesList;
  }

  static List<ActionOrReactionsParameter> storeActionReactionParameter(List<dynamic> parameter) {
    List<ActionOrReactionsParameter> parameterList = [];
    if (parameter.isEmpty) {
      return parameterList;
    }
    for (int i = 0; i < parameter.length; i++) {
      parameterList.add(
        ActionOrReactionsParameter.optional(
          parameter[i]["isexhaustive"] ?? true,
          parameter[i]["name"] ?? "",
          parameter[i]["route"] ?? "",
          parameter[i]["type"] ?? "",
          getValueListe(parameter[i]["values"]),
        )
      );
    }
    return parameterList;
  }

  static void storeServiceActionsOrReactionsData(dynamic responseData, String serviceName, String route) {
    List<ActionOrReactionList> listActionOrReaction = Globaldata.actionServicesActions;
    if (route == "reactions") {
      listActionOrReaction = Globaldata.reactionServicesReaction;
    }
    if (responseData is Map<String, dynamic> &&
        responseData[route] is List) {
      List<dynamic> actions = responseData[route];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          listActionOrReaction.add(ActionOrReactionList.optional(
            actionResponse['description'] ?? '',
            actionResponse['id'] ?? '',
            actionResponse['name'] ?? '',
            actionResponse['serviceid'] ?? '',
            serviceName,
            storeActionReactionParameter(actionResponse["parameters"] ?? []),
          ));
        }
      }
    }
  }

  static Future<void> storeServiceActionsOrReactions(String serviceName, bool hasActions, bool hasReactions) async {
    String getServiceActionRoute = "/actions";
    String getServiceReactionRoute = "/reactions";
    if (hasActions) {
      final response = await getServiceActionsOrReactions(serviceName + getServiceActionRoute);
      if (response == null) {
        return;
      }
      storeServiceActionsOrReactionsData(response, serviceName, "actions");
    }
    if (hasReactions) {
      final response = await getServiceActionsOrReactions(serviceName + getServiceReactionRoute);
      if (response == null) {
        return;
      }
      storeServiceActionsOrReactionsData(response, serviceName, "reactions");
    }
  }

  static void storeServiceData(List<String> servicesName, List<String> servicesColor,
      List<String> servicesIcons, List<String> servicesID, List<bool> serviceIsAuthNeeded, dynamic service) {
    servicesName.add(service['name']);
    servicesColor.add(service['color']);
    servicesIcons.add(service['logo']);
    serviceIsAuthNeeded.add(service['isauthneeded']);
    servicesID.add(service['id']);
  }

  static void setStoreServiceData(dynamic responseData) {
    final List<dynamic> responses = responseData;
    for (int i = 0; i < responses.length; i++) {
      final service = responseData[i];
      if (service is Map<String, dynamic>) {
        String actualServiceName = service['name'] ?? 'No Name';
        bool hasActions = service['hasactions'] ?? false;
        bool hasReactions = service['hasreactions'] ?? false;
        if (!hasActions && !hasReactions) {
          continue;
        }
        storeServiceActionsOrReactions(actualServiceName, hasActions, hasReactions);
        Globaldata.serviceList.add(
          ServiceData.optional(
            service['name'] ?? 'No Name',
            service['color'] ?? '',
            service['logo'] ?? '',
            service['id'] ?? '',
            service['description'] ?? '',
            service['isauthneeded'] ?? false,
            hasActions,
            hasReactions,
          ),
        );
      }
    }
    Future.delayed(const Duration(seconds: 1)).then((val) {
      LoginNavigator.redirectAndHandleGoBackRedirection(ExplorePage(), LoadServicePage());
    });
  }
}
import "package:http/src/response.dart";
import 'package:namer_app/globalData.dart';
import 'dart:convert';
import 'dart:developer' as developer;

import "../apiCall/apiRequest.dart";

class MyWorkflowsfunctions {
  static Future<dynamic> getWorkflowsResponse(String route) async {
    Response response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getWorkflowsResponse FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    return responseData;
  }

  static void storeServiceDescription(dynamic responseData, String serviceName) {
    if (responseData is Map<String, dynamic> &&
        responseData["actions"] is List) {
      List<dynamic> actions = responseData["actions"];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          Globaldata.actionServicesActions.add(ActionOrReactionList.optional(
            actionResponse['description'] ?? '',
            actionResponse['id'] ?? '',
            actionResponse['name'] ?? '',
            actionResponse['nbparam'] ?? 0,
            actionResponse['serviceid'] ?? '',
            serviceName,
          ));
        }
      }
    }
    if (responseData is Map<String, dynamic> &&
        responseData["reactions"] is List) {
      List<dynamic> actions = responseData["reactions"];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          Globaldata.reactionServicesReaction.add(ActionOrReactionList.optional(
            actionResponse['description'] ?? '',
            actionResponse['id'] ?? '',
            actionResponse['name'] ?? '',
            actionResponse['nbparam'] ?? 0,
            actionResponse['serviceid'] ?? '',
            serviceName,
          ));
        }
      }
    }
  }
}
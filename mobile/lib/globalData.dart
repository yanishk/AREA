import 'package:flutter/material.dart';
import 'login_page.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';

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

class GetApi {
  static Future<dynamic> get(String apiRoute) async {
    String apiUrl = dotenv.get('SERVER_LINK',
        fallback: 'Cannot find .env variable "SERVER_LINK"');
    String result = "";
    developer.log('SERVER_LINK : $apiUrl');
    try {
      developer.log('SERVER_LINK + redirection : ${apiUrl + apiRoute}');
      final response = await http.get(
        Uri.parse(apiUrl + apiRoute),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );
      final responseData = jsonDecode(response.body);
      return responseData;
    } catch (e) {
      developer.log("getServices ERROR: $e", level: 1);
      result = '$e';
      developer.log(result);
      if (result.split(' ').contains("SocketException:") == true ||
          (result.split(' ').contains("Invalid") == true &&
              result.split(' ').contains("argument(s):") == true)) {
        result = Globaldata.ErrorServerConnection;
      }
      return result;
    }
  }
}

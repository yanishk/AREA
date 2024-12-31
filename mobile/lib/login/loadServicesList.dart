import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:auto_size_text/auto_size_text.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

import "../custom_widget.dart";
import "../globalData.dart";
import "../explore/explorePage.dart";
import "../apiCall/apiRequest.dart";
import "loginNavigator.dart";

class LoadServicePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  LoadService createState() => LoadService();
}

class LoadService extends State<APICall> {
  String result = '';
  String apiRoute = 'actions-services';
  String title = "Select action service";
  String actionsServices = "actions-services";
  String reactionsServices = "reactions-services";
  bool isInit = true;
  bool hasToRetry = false;

  Future<dynamic> getWorkflowsResponse(String route) async {
    Response response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
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

  void storeServiceDescription(dynamic responseData, String serviceName, String route) {
    String typeToGet = "actions";
    List<ActionOrReactionList> listActionOrReaction = Globaldata.actionServicesActions;
    if (route != actionsServices) {
      typeToGet = "reactions";
      listActionOrReaction = Globaldata.reactionServicesReaction;
    }
    if (responseData is Map<String, dynamic> &&
        responseData[typeToGet] is List) {
      List<dynamic> actions = responseData[typeToGet];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          listActionOrReaction.add(ActionOrReactionList.optional(
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

  Future<void> getWorkflows(String route, String service) async {
    String getServiceActionOrReactionRoute = "/actions";
    if (route != actionsServices) {
      getServiceActionOrReactionRoute = "/reactions";
    }
    final response = await getWorkflowsResponse(service + getServiceActionOrReactionRoute);
    if (response == null) {
      return;
    }
    setState(() {
      storeServiceDescription(response, service, route);
    });
  }

  void storeData(List<String> servicesName, List<String> servicesColor,
      List<String> servicesIcons, List<String> servicesID, dynamic service) {
    servicesName.add(service['name']);
    servicesColor.add(service['color']);
    servicesIcons.add(service['logo']);
    servicesID.add(service['id']);
  }

  void setStoreData(dynamic responseData, String route) {
    final List<dynamic> allResponse = responseData;
    for (int i = 0; i < allResponse.length; i++) {
      final service = responseData[i];
      if (service is Map<String, dynamic>) {
        result = service['name'] ?? 'No Name';
        getWorkflows(route, result);
        if (route == actionsServices) {
          storeData(Globaldata.actionServicesName, Globaldata.actionServicesColor,
            Globaldata.actionServicesIcons, Globaldata.actionServicesID ,service);
        } else {
          storeData(
            Globaldata.reactionServicesName,
            Globaldata.reactionServicesColor,
            Globaldata.reactionServicesIcons,
            Globaldata.reactionServicesID,
            service
          );
          LoginNavigator.redirectAndHandleGoBackRedirection(ExplorePage(), LoadServicePage());
        }
      }
    }
  }

  Future<dynamic> getServiceResponse(String route) async {
    Response response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getServiceResponse FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    return responseData;
  }

  Future<void> getServices() async {
    final response = await getServiceResponse(actionsServices);
    final responseReaction = await getServiceResponse(reactionsServices);
    if (response == null || responseReaction == null) {
      hasToRetry = true;
      return;
    }
    setStoreData(response, actionsServices);
    setStoreData(responseReaction, reactionsServices);
  }

  @override
  Widget build(BuildContext context) {
    if (isInit == true) {
      isInit = false;
      getServices();
    }

    if (hasToRetry) {
      return Scaffold(
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding (
                padding: EdgeInsets.only (
                  top: MediaQuery.sizeOf(Globaldata.myContext).height / 3
                ),
              ),
              AutoSizeText(
                "Error during connection",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: const Color.fromARGB(255, 0, 0, 0)),
              ),
              Padding (padding: EdgeInsets.only (top: 50),),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText(
                  "Retry",
                  style: TextStyle(
                      fontSize: 30,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                ),
                onPressed: () {getServices();},
                size: Size(double.infinity, 55),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
              ),
            ],
          ),
        ),
      );
    } else {
      return Scaffold(
        body: Center(
          child: LoadingAnimationWidget.dotsTriangle(
            color: Color.fromARGB(255, 34, 34, 34),
            size: 200,
          ),
        ),
      );
    }
  }
}

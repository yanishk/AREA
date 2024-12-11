import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;

import "custom_widget.dart";
import "globalData.dart";
import "exploreAll.dart";
import 'apiCall/getApi.dart';

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
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  String result = '';
  String apiRoute = 'actions-services';
  String title = "Select action service";
  String actionsServices = "actions-services";
  String reactionsServices = "reactions-services";
  bool isInit = true;

  void storeData(List<String> servicesName, List<String> servicesColor,
      List<String> servicesIcons, dynamic service) {
    servicesName.add(service['Name']);
    servicesColor.add(service['Color']);
    servicesIcons.add(service['Logo']);
  }

  void setStoreData(dynamic responseData, String route) {
    final firstService = responseData[0];
    if (firstService is Map<String, dynamic>) {
      result = firstService['Name'] ?? 'No Name';
      if (route == actionsServices) {
        storeData(Globaldata.actionServicesName, Globaldata.actionServicesColor,
          Globaldata.actionServicesIcons, firstService);
      } else {
        storeData(
          Globaldata.reactionServicesName,
          Globaldata.reactionServicesColor,
          Globaldata.reactionServicesIcons,
          firstService
        );
        Navigator.push(
          Globaldata.myContext,
          MaterialPageRoute(builder: (context) => ExplorePageAll()),
        );
      }
    }
  }

  Future<dynamic> getServiceResponse(String route) async {
    Response response = await GetApi.get(route);
    if (response.statusCode == 504) {
      setState(() {
        developer.log('${response.body}');
      });
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getServiceResponse FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    developer.log('pass');
    return responseData;
  }

  Future<void> getServices() async {

    final response = await getServiceResponse(actionsServices);
    if (response == null) {
      return;
    }
    setStoreData(response, actionsServices);
    final responseReaction = await getServiceResponse(reactionsServices);
    if (responseReaction == null) {
      return;
    }
    setStoreData(responseReaction, reactionsServices);
  }

  @override
  Widget build(BuildContext context) {
    if (isInit == true) {
      isInit = false;
      getServices();
    }

    return Column(
      children: const [
        MyDivider(),
      ],
    );
  }
}

import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:auto_size_text/auto_size_text.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

import "../../custom_widget.dart";
import "../../globalData.dart";
import "../../apiCall/apiRequest.dart";
import 'loadDataFunctions.dart';

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
  bool hasToRetry = false;

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
    final allServicesResponses = await getServiceResponse("services");
    if (allServicesResponses == null) {
      hasToRetry = true;
      return;
    }
    LoadDataFunctions.setStoreServiceData(allServicesResponses);
  }

  @override
  void initState() {
    super.initState();
    Globaldata.resetActionReactionServicesList();
    LoadDataFunctions.storeUserInfos();
    getServices();
  }

  @override
  Widget build(BuildContext context) {
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

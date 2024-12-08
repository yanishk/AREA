import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';

import "custom_widget.dart";
import "globalData.dart";
import "exploreAll.dart";

class LoadServicePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  LoadService createState() => LoadService();
}

class LoadService extends State<APICall> {
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  String result = '';
  String apiRoute = 'actions-services';
  String title = "Select action service";
  String actionsServices = "actions-services";
  String reactionsServices = "reactions-services";
  bool isInit = true;

  Future<void> getServices() async {
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
      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body);

        setState(() {
          developer.log("getServices OK :\n$responseData");
          if (responseData is List && responseData.isNotEmpty) {
            final firstService = responseData[0];
            if (firstService is Map<String, dynamic>) {
              result = firstService['Name'] ?? 'No Name';
              developer.log("getServices OK :\n$result");
              if (apiRoute == actionsServices) {
                Globaldata.actionServicesName.add(firstService['Name']);
                Globaldata.actionServicesColor.add(firstService['Color']);
                Globaldata.actionServicesIcons.add(firstService['Logo']);
                developer.log("getServices $actionsServices");
                apiRoute = reactionsServices;
                getServices();
              } else {
                Globaldata.reactionServicesName.add(firstService['Name']);
                Globaldata.reactionServicesColor.add(firstService['Color']);
                Globaldata.reactionServicesIcons.add(firstService['Logo']);
                developer.log("getServices $reactionsServices");
                Navigator.push(
                  Globaldata.myContext,
                  MaterialPageRoute(builder: (context) => ExplorePageAll()),
                );
              }
            }
          } else {
            result = 'No data available';
          }
        });
      } else {
        final responseData = jsonDecode(response.body);
        developer.log("getServices FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
        setState(() {
          result = '${responseData['error']}';
        });
      }
    } catch (e) {
      developer.log("getServices ERROR: $e", level: 1);
      setState(() {
        result = '$e';
        developer.log(result);
        if (result.split(' ').contains("SocketException:") == true ||
           (result.split(' ').contains("Invalid") == true && result.split(' ').contains("argument(s):") == true)) {
          result = Globaldata.ErrorServerConnection;
        }
      });
    }
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

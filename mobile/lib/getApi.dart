import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';

import "custom_widget.dart";
import "globalData.dart";

class GetApiPage extends StatelessWidget {
  const GetApiPage({
    super.key,
    required this.apiRoute,
  });
  final String apiRoute;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: GetAPICall(apiRoute: apiRoute),
    );
  }
}

class GetAPICall extends StatefulWidget {
  const GetAPICall({
    super.key,
    required this.apiRoute,
  });
  final String apiRoute;
  @override
  GetApi createState() => GetApi(apiRoute: apiRoute);
}

class GetApi extends State<GetAPICall> {
  GetApi({
    required this.apiRoute,
  });
  final String apiRoute;
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  String result = '';
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
          Navigator.pop(Globaldata.myContext, responseData);
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

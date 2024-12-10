import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:fluttertoast/fluttertoast.dart';

import "custom_widget.dart";
import "globalData.dart";
import "exploreAll.dart";

class SaveNewWorkflowsPage extends StatelessWidget {
  const SaveNewWorkflowsPage({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(isAnAction: isAnAction),
    );
  }
}

class APICall extends StatefulWidget {
  const APICall({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  SaveNewWorkflows createState() => SaveNewWorkflows(isAnAction: isAnAction);
}

class SaveNewWorkflows extends State<APICall> {
  SaveNewWorkflows({
    required this.isAnAction,
  });
  final bool isAnAction;
  TextEditingController workflowTitleController = new TextEditingController(
      text:
          "If ${Globaldata.serviceAction.name} ${Globaldata.serviceAction.parameter}, Then ${Globaldata.reactionServicesName[Globaldata.secondService]} : ${Globaldata.serviceReaction.name}");
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  final String apiRoute = "workflow/create";
  String result = '';
  String email = "";
  bool isInit = true;

  Future<void> postCreateWorkflows() async {
    developer.log('SERVER_LINK : $apiUrl');
    try {
      developer.log('SERVER_LINK + redirection : ${apiUrl + apiRoute}');
      final response = await http.post(
        Uri.parse('${apiUrl + apiRoute}'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
        body: jsonEncode(<String, dynamic>{
          'actionid': Globaldata.serviceAction.id,
          'actionparam': Globaldata.serviceAction.parameterList,
          'name': workflowTitleController.value.text.toString(),
          'reactionid': Globaldata.serviceReaction.id,
          'reactionparam': Globaldata.serviceReaction.parameterList,
        }),
      );

      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body);
        setState(() {
          result = responseData.toString();
          Fluttertoast.showToast(
              msg: "Workflows successfully created",
              toastLength: Toast.LENGTH_SHORT,
              timeInSecForIosWeb: 1,
              backgroundColor: const Color.fromARGB(255, 70, 70, 70),
              textColor: Colors.white,
              fontSize: 24.0);
          Navigator.push(
            Globaldata.myContext,
            MaterialPageRoute(builder: (context) => ExplorePageAll()),
          );
        });
      } else {
        final responseData = jsonDecode(response.body);
        setState(() {
          result = '${responseData['error']}';
          developer.log("Failed : $result");
        });
      }
    } catch (e) {
      setState(() {
        result = '$e';
        developer.log(result);
        if (result.split(' ').contains("SocketException:") == true ||
            (result.split(' ').contains("Invalid") == true &&
                result.split(' ').contains("argument(s):") == true)) {
          result = Globaldata.ErrorServerConnection;
        }
      });
    }
  }

  Future<void> getEmail() async {
    final responce = await GetApi.get("user/email");
    email = responce["email"];
  }

  @override
  Widget build(BuildContext context) {
    if (isInit == true) {
      isInit = false;
      getEmail();
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Container(
              decoration: BoxDecoration(
                color: Color.fromARGB(255, 63, 68, 69),
                border: Border.all(
                  color: Color.fromARGB(255, 63, 68, 69),
                  width: 2.0,
                ),
                borderRadius: BorderRadius.circular(4.0),
              ),
              child: InkWell(
                child: Column(
                  children: [
                    Padding(padding: EdgeInsets.only(top: 30)),
                    SelectServiceTopBar(
                        title: "Review and finish", color: Colors.white),
                    Padding(padding: EdgeInsets.only(top: 20)),
                    LeftWidget(
                        child: Row(
                      children: [
                        SafeWebImage(
                          url:
                              "${Globaldata.domainName}${Globaldata.actionServicesIcons[Globaldata.firstService]}",
                          width: 30,
                          height: 30,
                        ),
                        Padding(padding: EdgeInsets.only(left: 10)),
                        SafeWebImage(
                          url:
                              "${Globaldata.domainName}${Globaldata.reactionServicesIcons[Globaldata.secondService]}",
                          width: 30,
                          height: 30,
                        ),
                      ],
                    )),
                    Padding(padding: EdgeInsets.only(top: 5)),
                    LeftWidget(
                      child: AutoSizeText(
                        "Workflow title",
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                        maxLines: 1,
                      ),
                    ),
                    WidgetCard(
                      backgroundColor: Colors.white,
                      borderRadius: 10,
                      child: CustomClassicForm(
                        controller: workflowTitleController,
                        maxLines: 5,
                        hint: 'Email',
                        padding:
                            EdgeInsets.symmetric(horizontal: 0, vertical: 0),
                        focusedBorderColor: Colors.white,
                        borderWidth: 0,
                      ),
                    ),
                    LeftWidget(
                      child: AutoSizeText(
                        "By $email",
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                        maxLines: 1,
                      ),
                    ),
                    Padding(padding: EdgeInsets.only(top: 10)),
                  ],
                ),
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Finish",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: () {
                postCreateWorkflows();
              },
              size: Size(double.infinity, 60),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(
              Icons.arrow_back,
              size: 32,
              color: color,
            ),
          ),
        ),
        Center(
          child: AutoSizeText(
            title,
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: color,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}

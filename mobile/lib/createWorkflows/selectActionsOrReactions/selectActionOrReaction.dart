import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../../custom_widget.dart";
import "../../globalData.dart";
import '../createWidget.dart';
import '../createWorkflowsData.dart';
import "selectActionOrReactionParam.dart";
import 'actionOrReactionsFunctions.dart';

class SelectActionOrReactionPage extends StatelessWidget {
  const SelectActionOrReactionPage({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: SelectActionOrReactionAPICall(serviceIndex: serviceIndex),
    );
  }
}

class SelectActionOrReactionAPICall extends StatefulWidget {
  const SelectActionOrReactionAPICall({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  SelectActionOrReaction createState() => SelectActionOrReaction(serviceIndex: serviceIndex);
}

class SelectActionOrReaction extends State<SelectActionOrReactionAPICall> {
  SelectActionOrReaction({
    required this.serviceIndex,
  });
  final int serviceIndex;
  List<ActionOrReactionData> actionOrReactionList = [];
  Color colorTheme = Colors.black;
  String apiRoute = "actions";
  String actionDescription = "";
  List<Widget> ActionOrReactionButtonsList = [];
  Size buttonSize = Size(MediaQuery.sizeOf(Globaldata.myContext).width, 20);
  double buttonRadius = 8;

  @override
  void initState() {
    super.initState();
    if (!isAnAction) {
      setState(() {
        apiRoute = "reactions";
        CreateData.selectedActionOrReactionData = ActionOrReactionData.optional();
        CreateData.selectedServiceId = serviceIndex;
      });
    }
    CreateData.hasToLogin = false;
    actionDescription = CreateData.actualServiceList[serviceIndex].description;
    getActionsOrReactions();
    isAlreadyLogged();
  }

  Future<void> isAlreadyLogged() async {
    setState(() {
      CreateData.hasToLogin = CreateData.actualServiceList[serviceIndex].isAuthNeeded;
    });
    if (CreateData.actualServiceList[serviceIndex].isAuthNeeded) {
      final bool isLogged = await ActionsOrReactionsFunctions.isAlreadyLogged(CreateData.actualServiceList[serviceIndex]);
      CreateData.hasToLogin = !isLogged;
    }
  }

  Future<void> getActionsOrReactions() async {
    setState(() {
      if (apiRoute == "actions") {
        storeServiceDescription(Globaldata.actionServicesActions);
      } else {
        storeServiceDescription(Globaldata.reactionServicesReaction);
      }
    });
  }

  void storeServiceDescription(List<ActionOrReactionList> allActionOrReactionList) {
    for (int i = 0; i < allActionOrReactionList.length; i++) {
      if (allActionOrReactionList[i].servicename == CreateData.actualServiceList[serviceIndex].name) {
        actionOrReactionList.add(ActionOrReactionData.optional(
          allActionOrReactionList[i].id,
          allActionOrReactionList[i].name,
          allActionOrReactionList[i].description,
          allActionOrReactionList[i].serviceid,
          CreateData.actualServiceList[serviceIndex].icon,
          allActionOrReactionList[i].parameter,
        ));
        ActionOrReactionButtonsList.add(
          actionOrReactionButtons(actionOrReactionList.length - 1),
        );
      }
    }
  }

  Widget actionOrReactionButtons(int index) {
    return Column(children: [
      Padding(padding: EdgeInsets.only(top: 15)),
      CustomButton(
        backgroundColor: HexColor(CreateData.actualServiceList[serviceIndex].color),
        text: ButtonText(text: actionOrReactionList[index].name,),
        onPressed: () {
          saveSelectedActionOrReaction(index);
        },
        size: buttonSize,
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        borderCircularRadius: buttonRadius,
      )
    ]);
  }

  Future<void> saveSelectedActionOrReaction(int index) async {
    CreateData.selectedServiceId = serviceIndex;
    CreateData.selectedActionOrReactionData = actionOrReactionList[index];
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => SelectParamActionOrReactionPage()),
    );
    if (!CreateData.selectOrUpdateWorkflow) {
      Navigator.pop(Globaldata.myContext);
      return;
    }
    actionOrReactionList[index] = CreateData.selectedActionOrReactionData;
    if (isAnAction) {
      CreateData.serviceAction = actionOrReactionList[index];
    } else {
      CreateData.serviceReaction = actionOrReactionList[index];
    }
    Navigator.pop(Globaldata.myContext);
  }

  @override
  Widget build(BuildContext context) {
    colorTheme = HexColor(CreateData.actualServiceList[serviceIndex].color);
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            ActionOrReactionDescription(
                colorTheme: colorTheme,
                title: "Select $apiRoute",
                serviceIcon: CreateData.actualServiceList[serviceIndex].icon,
                serviceName: CreateData.actualServiceList[serviceIndex].name,
                actionDescription: actionDescription),
            Padding(padding: EdgeInsets.only(top: 15)),
            Column(
              children: ActionOrReactionButtonsList,
            ),
            Padding(padding: EdgeInsets.only(top: 15)),
          ],
        ),
      ),
    );
  }
}

class ButtonText extends StatelessWidget {
  const ButtonText({
    super.key,
    required this.text,
  });
  final String text;

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: Alignment.centerLeft,
      child: new AutoSizeText(
        text,
        style: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.bold,
          color: Colors.white,
          height: 2.3,
        ),
      )
    );
  }
}

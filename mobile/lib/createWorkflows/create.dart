import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import "package:namer_app/custom_widget.dart";
import "package:namer_app/footer_page.dart";
import "package:namer_app/createWorkflows/selectService.dart";

import "../globalData.dart";
import "saveNewWorkflows.dart";
import "createWorkflowsData.dart";

class CreateWorkflowsPage extends StatelessWidget {
  const CreateWorkflowsPage({super.key});

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
  const APICall({super.key});

  @override
  CreateWorkflows createState() => CreateWorkflows();
}

class CreateWorkflows extends State<APICall> {
  bool firstButtonComplete = false;
  int secondButtonState = 0;
  Color enableButton = Color.fromARGB(255, 34, 34, 34);
  Color firstButtonColor = Color.fromARGB(255, 34, 34, 34);
  Color secondButtonColor = Color.fromARGB(255, 150, 150, 150);
  double firstButtonHight = 70;
  String firstSelectedService = "";
  String secondSelectedService = "";

  Widget firstButton() {
    if (firstButtonComplete) {
      return CreateSelectedButtonContent(
        firstText: "If ",
        icon: CreateData.serviceAction.icon,
        actionOrReactionName: CreateData.serviceAction.name,
      );
    }
    return PlaceableWidget(
      left: MediaQuery.sizeOf(Globaldata.myContext).width / 10,
      child: DoubleWidget(
        firstChild: AutoSizeText(
          "If This",
          style: TextStyle(
              fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
        ),
        secondChild: TextWithBackground(text: "Add", fontSize: 25),
        padding: EdgeInsets.only(
            left: MediaQuery.sizeOf(Globaldata.myContext).width / 5),
      ),
    );
  }

  Widget secondButton() {
    if (secondButtonState == 1) {
      return DoubleWidget(
        firstChild: PlaceableWidget(
          child: AutoSizeText(
            "Then That",
            style: TextStyle(
                fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
          ),
        ),
        secondChild: TextWithBackground(text: "Add", fontSize: 25),
        padding: EdgeInsets.only(
            left: MediaQuery.sizeOf(Globaldata.myContext).width / 10),
      );
    } else if (secondButtonState == 2) {
      return CreateSelectedButtonContent(
        firstText: "Then ",
        icon: CreateData.serviceReaction.icon,
        actionOrReactionName: CreateData.serviceReaction.name,
      );
    }
    return const AutoSizeText(
      "Then That",
      style: TextStyle(
          fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
    );
  }

  Widget continueButton() {
    if (secondButtonState == 2) {
      return CustomButton(
        backgroundColor: Color.fromARGB(255, 34, 34, 34),
        shadowSize: 0,
        text: AutoSizeText(
          "Continue",
          style: TextStyle(
              fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
        ),
        onPressed: () {
          Navigator.push(Globaldata.myContext,
              MaterialPageRoute(builder: (context) => SaveNewWorkflowsPage()));
        },
        size: Size(double.infinity, 60),
        padding: EdgeInsets.symmetric(horizontal: 15.0),
      );
    }
    return SizedBox();
  }

  Future<void> firstButtonNextState() async {
    setState(() {
      isAnAction = true;
    });
    final result = await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(
          builder: (context) => SelectServicePage()),
    );
    setState(() {
      if (firstButtonComplete == false && result != null) {
        firstButtonComplete = true;
        secondButtonState++;
        firstButtonColor = HexColor(CreateData.actionServiceList[Globaldata.actionServiceIndex].color);
        secondButtonColor = enableButton;
      }
      if (result != null) {
        firstSelectedService = "$result";
        firstButtonColor = HexColor(CreateData.actionServiceList[Globaldata.actionServiceIndex].color);
      }
    });
  }

  Future<void> secondButtonNextState() async {
    if (secondButtonState != 0) {
      isAnAction = false;
      final result = await Navigator.push(
        Globaldata.myContext,
        MaterialPageRoute(
            builder: (context) => SelectServicePage()),
      );
      setState(() {
        if (result != null) {
          secondButtonColor = HexColor(CreateData.reactionServiceList[Globaldata.reactionServiceIndex].color);
          secondButtonState = 2;
          secondSelectedService = "$result";
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            const Padding(padding: EdgeInsets.only(top: 46)),
            const AutoSizeText(
              "Create",
              style: TextStyle(
                fontSize: 30,
                fontWeight: FontWeight.w900,
                color: Colors.black,
              ),
              maxLines: 1,
            ),
            const Padding(padding: EdgeInsets.only(top: 30)),
            CustomButton(
              backgroundColor: firstButtonColor,
              shadowSize: 0,
              text: firstButton(),
              onPressed: () {
                firstButtonNextState();
              },
              size: Size(double.infinity, firstButtonHight),
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              borderCircularRadius: 8,
            ),
            const MyVerticalDivider(
              height: 40,
              width: 5,
              color: Color.fromARGB(255, 222, 222, 222),
            ),
            CustomButton(
              backgroundColor: secondButtonColor,
              shadowSize: 0,
              text: secondButton(),
              onPressed: () {
                secondButtonNextState();
              },
              size: const Size(double.infinity, 70),
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              borderCircularRadius: 8,
            ),
            Padding(padding: EdgeInsets.only(top: 50)),
            continueButton(),
            Padding(padding: EdgeInsets.only(top: 80)),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 2),
    );
  }
}

class CreateSelectedButtonContent extends StatelessWidget {
  const CreateSelectedButtonContent({
    super.key,
    required this.firstText,
    required this.icon,
    required this.actionOrReactionName,
  });
  final String firstText;
  final String icon;
  final String actionOrReactionName;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        AutoSizeText(
          firstText,
          style: TextStyle (
            fontSize: 30,
            fontWeight: FontWeight.bold,
            color: Colors.white
          ),
        ),
        SafeWebImage(
          url: "${Globaldata.domainName}$icon",
          width: 30,
          height: 30,
        ),
        Padding(padding: EdgeInsets.only(right: 10)),
        Expanded(
          child: AutoSizeText(
            " $actionOrReactionName",
            style: TextStyle (
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white
            ),
          ),
        )
      ]
    );
  }
}
